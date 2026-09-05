import{$r as e,An as t,Fr as n,Ii as r,Jn as i,Jr as a,Qr as o,Ri as s,Xr as c,Yr as l,Zr as u,ar as d,ei as f,ki as p,qn as m,ti as h,wr as g}from"./router-hzdBzcXs.js";import{n as _}from"./Space-CPehM0jl.js";import{s as v}from"./index-DiTfxJos.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=a([l(`descriptions`,{fontSize:`var(--n-font-size)`},[l(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[l(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),l(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),o(`bordered`,[l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[a(`&:last-child`,[l(`descriptions-table-content`,{paddingBottom:0})])])])])]),u(`left-label-placement`,[l(`descriptions-table-content`,[a(`> *`,{verticalAlign:`top`})])]),u(`left-label-align`,[a(`th`,{textAlign:`left`})]),u(`center-label-align`,[a(`th`,{textAlign:`center`})]),u(`right-label-align`,[a(`th`,{textAlign:`right`})]),u(`bordered`,[l(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[l(`descriptions-table`,[l(`descriptions-table-row`,[a(`&:not(:last-child)`,[l(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),l(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[a(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-content`,[a(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),l(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),l(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[l(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[l(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[l(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),l(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[c(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),c(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),f(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),h(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=r({name:`Descriptions`,props:Object.assign(Object.assign({},t.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(n){let{mergedClsPrefixRef:r,inlineThemeDisabled:a,mergedComponentPropsRef:o}=i(n),s=p(()=>n.size||o?.value?.Descriptions?.size||`medium`),c=t(`Descriptions`,`-descriptions`,b,v,n,r),l=p(()=>{let{bordered:t}=n,r=s.value,{common:{cubicBezierEaseInOut:i},self:{titleTextColor:a,thColor:o,thColorModal:l,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[e(`fontSize`,r)]:S,[e(t?`thPaddingBordered`:`thPadding`,r)]:C,[e(t?`tdPaddingBordered`:`tdPadding`,r)]:w}}=c.value;return{"--n-title-text-color":a,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":i,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":o,"--n-th-color-modal":l,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),u=a?m(`descriptions`,p(()=>{let e=``,{bordered:t}=n;return t&&(e+=`a`),e+=s.value[0],e}),l,n):void 0;return{mergedClsPrefix:r,cssVars:a?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender,compitableColumn:g(n,[`columns`,`column`]),inlineThemeDisabled:a,mergedSize:s}},render(){let e=this.$slots.default,t=e?d(e()):[];t.length;let{contentClass:r,labelClass:i,compitableColumn:a,labelPlacement:o,labelAlign:c,mergedSize:l,bordered:u,title:f,cssVars:p,mergedClsPrefix:m,separator:h,onRender:g}=this;g?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let c=t.props||{},l=v.length-1===n,d=[`label`in c?c.label:y(t,`label`)],f=[y(t)],p=c.span||1,g=e.span;e.span+=p;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(o===`left`)u?e.row.push(s(`th`,{class:[`${m}-descriptions-table-header`,i],colspan:1,style:_},d),s(`td`,{class:[`${m}-descriptions-table-content`,r],colspan:l?(a-g)*2+1:p*2-1,style:b},f)):e.row.push(s(`td`,{class:`${m}-descriptions-table-content`,colspan:l?(a-g)*2:p*2},s(`span`,{class:[`${m}-descriptions-table-content__label`,i],style:_},[...d,h&&s(`span`,{class:`${m}-descriptions-separator`},h)]),s(`span`,{class:[`${m}-descriptions-table-content__content`,r],style:b},f)));else{let t=l?(a-g)*2:p*2;e.row.push(s(`th`,{class:[`${m}-descriptions-table-header`,i],colspan:t,style:_},d)),e.secondRow.push(s(`td`,{class:[`${m}-descriptions-table-content`,r],colspan:t,style:b},f))}return(e.span>=a||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),o!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>s(`tr`,{class:`${m}-descriptions-table-row`},e));return s(`div`,{style:p,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${o}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,u&&`${m}-descriptions--bordered`]},f||this.$slots.header?s(`div`,{class:`${m}-descriptions-header`},f||_(this,`header`)):null,s(`div`,{class:`${m}-descriptions-table-wrapper`},s(`table`,{class:`${m}-descriptions-table`},s(`tbody`,null,o===`top`&&s(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},n(a*2,s(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=r({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};