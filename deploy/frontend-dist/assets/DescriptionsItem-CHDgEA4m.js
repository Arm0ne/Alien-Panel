import{$n as e,$r as t,Br as n,Fn as r,Pi as i,Qn as a,Ui as o,Vi as s,ai as c,ei as l,ii as u,kr as d,ni as f,oi as p,ri as m,ti as h,ur as g}from"./router-B5zIVf8_.js";import{n as _}from"./Space-Ds8-hval.js";import{s as v}from"./index-DTohGTz4.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=t([l(`descriptions`,{fontSize:`var(--n-font-size)`},[l(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[l(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),l(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),m(`bordered`,[l(`descriptions-table-wrapper`,[l(`descriptions-table`,[l(`descriptions-table-row`,[t(`&:last-child`,[l(`descriptions-table-content`,{paddingBottom:0})])])])])]),f(`left-label-placement`,[l(`descriptions-table-content`,[t(`> *`,{verticalAlign:`top`})])]),f(`left-label-align`,[t(`th`,{textAlign:`left`})]),f(`center-label-align`,[t(`th`,{textAlign:`center`})]),f(`right-label-align`,[t(`th`,{textAlign:`right`})]),f(`bordered`,[l(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[l(`descriptions-table`,[l(`descriptions-table-row`,[t(`&:not(:last-child)`,[l(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),l(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[t(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),l(`descriptions-table-content`,[t(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),l(`descriptions-header`,`
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
 `,[h(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),h(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),c(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),p(l(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=s({name:`Descriptions`,props:Object.assign(Object.assign({},r.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(t){let{mergedClsPrefixRef:n,inlineThemeDisabled:o,mergedComponentPropsRef:s}=e(t),c=i(()=>t.size||s?.value?.Descriptions?.size||`medium`),l=r(`Descriptions`,`-descriptions`,b,v,t,n),f=i(()=>{let{bordered:e}=t,n=c.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:s,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[u(`fontSize`,n)]:S,[u(e?`thPaddingBordered`:`thPadding`,n)]:C,[u(e?`tdPaddingBordered`:`tdPadding`,n)]:w}}=l.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":s,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),p=o?a(`descriptions`,i(()=>{let e=``,{bordered:n}=t;return n&&(e+=`a`),e+=c.value[0],e}),f,t):void 0;return{mergedClsPrefix:n,cssVars:o?void 0:f,themeClass:p?.themeClass,onRender:p?.onRender,compitableColumn:d(t,[`columns`,`column`]),inlineThemeDisabled:o,mergedSize:c}},render(){let e=this.$slots.default,t=e?g(e()):[];t.length;let{contentClass:r,labelClass:i,compitableColumn:a,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:u,title:d,cssVars:f,mergedClsPrefix:p,separator:m,onRender:h}=this;h?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let c=t.props||{},l=v.length-1===n,d=[`label`in c?c.label:y(t,`label`)],f=[y(t)],h=c.span||1,g=e.span;e.span+=h;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(s===`left`)u?e.row.push(o(`th`,{class:[`${p}-descriptions-table-header`,i],colspan:1,style:_},d),o(`td`,{class:[`${p}-descriptions-table-content`,r],colspan:l?(a-g)*2+1:h*2-1,style:b},f)):e.row.push(o(`td`,{class:`${p}-descriptions-table-content`,colspan:l?(a-g)*2:h*2},o(`span`,{class:[`${p}-descriptions-table-content__label`,i],style:_},[...d,m&&o(`span`,{class:`${p}-descriptions-separator`},m)]),o(`span`,{class:[`${p}-descriptions-table-content__content`,r],style:b},f)));else{let t=l?(a-g)*2:h*2;e.row.push(o(`th`,{class:[`${p}-descriptions-table-header`,i],colspan:t,style:_},d)),e.secondRow.push(o(`td`,{class:[`${p}-descriptions-table-content`,r],colspan:t,style:b},f))}return(e.span>=a||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>o(`tr`,{class:`${p}-descriptions-table-row`},e));return o(`div`,{style:f,class:[`${p}-descriptions`,this.themeClass,`${p}-descriptions--${s}-label-placement`,`${p}-descriptions--${c}-label-align`,`${p}-descriptions--${l}-size`,u&&`${p}-descriptions--bordered`]},d||this.$slots.header?o(`div`,{class:`${p}-descriptions-header`},d||_(this,`header`)):null,o(`div`,{class:`${p}-descriptions-table-wrapper`},o(`table`,{class:`${p}-descriptions-table`},o(`tbody`,null,s===`top`&&o(`tr`,{class:`${p}-descriptions-table-row`,style:{visibility:`collapse`}},n(a*2,o(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=s({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};