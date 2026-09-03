import{$r as e,Di as t,Gn as n,Ii as r,Jr as i,Kn as a,Kr as o,Nr as s,On as c,Pi as l,Qr as u,Sr as d,Xr as f,Yr as p,Zr as m,qr as h,rr as g}from"./router-B7Ig3Ii3.js";import{n as _}from"./Space-B_qiysju.js";import{s as v}from"./index-Cvva8JXN.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=o([h(`descriptions`,{fontSize:`var(--n-font-size)`},[h(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),h(`descriptions-table-wrapper`,[h(`descriptions-table`,[h(`descriptions-table-row`,[h(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),h(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),f(`bordered`,[h(`descriptions-table-wrapper`,[h(`descriptions-table`,[h(`descriptions-table-row`,[o(`&:last-child`,[h(`descriptions-table-content`,{paddingBottom:0})])])])])]),p(`left-label-placement`,[h(`descriptions-table-content`,[o(`> *`,{verticalAlign:`top`})])]),p(`left-label-align`,[o(`th`,{textAlign:`left`})]),p(`center-label-align`,[o(`th`,{textAlign:`center`})]),p(`right-label-align`,[o(`th`,{textAlign:`right`})]),p(`bordered`,[h(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[h(`descriptions-table`,[h(`descriptions-table-row`,[o(`&:not(:last-child)`,[h(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),h(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),h(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[o(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),h(`descriptions-table-content`,[o(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),h(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),h(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[h(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[h(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[h(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),h(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[i(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),i(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),u(h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),e(h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=l({name:`Descriptions`,props:Object.assign(Object.assign({},c.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:o}=a(e),s=t(()=>e.size||o?.value?.Descriptions?.size||`medium`),l=c(`Descriptions`,`-descriptions`,b,v,e,r),u=t(()=>{let{bordered:t}=e,n=s.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:c,thTextColor:u,thFontWeight:d,tdTextColor:f,tdColor:p,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[m(`fontSize`,n)]:S,[m(t?`thPaddingBordered`:`thPadding`,n)]:C,[m(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=l.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":u,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":c,"--n-td-color":p,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),f=i?n(`descriptions`,t(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=s.value[0],t}),u,e):void 0;return{mergedClsPrefix:r,cssVars:i?void 0:u,themeClass:f?.themeClass,onRender:f?.onRender,compitableColumn:d(e,[`columns`,`column`]),inlineThemeDisabled:i,mergedSize:s}},render(){let e=this.$slots.default,t=e?g(e()):[];t.length;let{contentClass:n,labelClass:i,compitableColumn:a,labelPlacement:o,labelAlign:c,mergedSize:l,bordered:u,title:d,cssVars:f,mergedClsPrefix:p,separator:m,onRender:h}=this;h?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,s)=>{let c=t.props||{},l=v.length-1===s,d=[`label`in c?c.label:y(t,`label`)],f=[y(t)],h=c.span||1,g=e.span;e.span+=h;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(o===`left`)u?e.row.push(r(`th`,{class:[`${p}-descriptions-table-header`,i],colspan:1,style:_},d),r(`td`,{class:[`${p}-descriptions-table-content`,n],colspan:l?(a-g)*2+1:h*2-1,style:b},f)):e.row.push(r(`td`,{class:`${p}-descriptions-table-content`,colspan:l?(a-g)*2:h*2},r(`span`,{class:[`${p}-descriptions-table-content__label`,i],style:_},[...d,m&&r(`span`,{class:`${p}-descriptions-separator`},m)]),r(`span`,{class:[`${p}-descriptions-table-content__content`,n],style:b},f)));else{let t=l?(a-g)*2:h*2;e.row.push(r(`th`,{class:[`${p}-descriptions-table-header`,i],colspan:t,style:_},d)),e.secondRow.push(r(`td`,{class:[`${p}-descriptions-table-content`,n],colspan:t,style:b},f))}return(e.span>=a||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),o!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>r(`tr`,{class:`${p}-descriptions-table-row`},e));return r(`div`,{style:f,class:[`${p}-descriptions`,this.themeClass,`${p}-descriptions--${o}-label-placement`,`${p}-descriptions--${c}-label-align`,`${p}-descriptions--${l}-size`,u&&`${p}-descriptions--bordered`]},d||this.$slots.header?r(`div`,{class:`${p}-descriptions-header`},d||_(this,`header`)):null,r(`div`,{class:`${p}-descriptions-table-wrapper`},r(`table`,{class:`${p}-descriptions-table`},r(`tbody`,null,o===`top`&&r(`tr`,{class:`${p}-descriptions-table-row`,style:{visibility:`collapse`}},s(a*2,r(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=l({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};